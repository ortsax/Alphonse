package plugins

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
)

// sendTimeout caps the server-ACK wait per send.
// 20 s is well above normal WhatsApp RTT but prevents a single stuck send
// from holding messageSendLock for the default 75 s.
const sendTimeout = 20 * time.Second

// maxConcurrentSends is the maximum number of messages that may be in the
// encrypt+write+ACKwait phase simultaneously.  whatsmeow's messageSendLock
// already serialises the encrypt+write step, so this bound only limits how
// many ACK waits can overlap at once.  8 provides good burst throughput while
// staying well within WhatsApp's server-side flow-control limits.
const maxConcurrentSends = 8

type sendTask struct {
	client    *whatsmeow.Client
	to        types.JID
	msg       *waProto.Message
	id        types.MessageID
	enqueuedAt time.Time // set by Queue helpers for latency tracking
}

// sendQueue buffers fire-and-forget outgoing messages.
// Capacity 512 ensures a burst of concurrent commands never blocks a handler.
var sendQueue = make(chan sendTask, 512)

// sendSem limits how many sends may be in-flight at once.
var sendSem = make(chan struct{}, maxConcurrentSends)

func init() {
	go sendWorker()
}

// sendWorker drains sendQueue using a bounded goroutine pool.
//
// whatsmeow's messageSendLock serialises the encrypt+write step so Signal
// session ordering is always correct.  With the early-unlock patch in
// patched/send.go the lock is released as soon as the frame is on the wire,
// letting the next goroutine begin its own encrypt+write while the previous
// one awaits the server ACK.  The semaphore caps concurrent ACK waits so we
// don't flood the server.
func sendWorker() {
	for task := range sendQueue {
		sendSem <- struct{}{} // acquire slot (blocks if maxConcurrentSends in flight)
		go func(t sendTask) {
			defer func() { <-sendSem }()

			queueWait := time.Since(t.enqueuedAt)
			start := time.Now()

			resp, err := t.client.SendMessage(
				context.Background(),
				t.to,
				t.msg,
				whatsmeow.SendRequestExtra{
					ID:      t.id,
					Timeout: sendTimeout,
				},
			)
			if err != nil {
				fmt.Fprintf(os.Stderr, "[Send ERROR] %s → %s: %v\n", t.id, t.to, err)
			}

			total := time.Since(start)
			dt := resp.DebugTimings
			// onWire = time until message was written to socket (user-visible send latency)
			// ack    = server confirmation wait (async, does NOT block next send)
			onWire := dt.Queue + dt.LIDFetch + dt.Marshal + dt.GetParticipants + dt.GetDevices + dt.GroupEncrypt + dt.PeerEncrypt + dt.Send
			fmt.Fprintf(os.Stderr,
				"[Send TIMING] id=%s to=%s  queue_wait=%s on_wire=%s ack=%s total=%s | lock=%s lid=%s enc=%s write=%s retry=%s\n",
				t.id, t.to,
				queueWait.Round(time.Millisecond),
				onWire.Round(time.Millisecond),
				dt.Resp.Round(time.Millisecond),
				total.Round(time.Millisecond),
				dt.Queue.Round(time.Millisecond),
				dt.LIDFetch.Round(time.Millisecond),
				(dt.Marshal+dt.GetParticipants+dt.GetDevices+dt.GroupEncrypt+dt.PeerEncrypt).Round(time.Millisecond),
				dt.Send.Round(time.Millisecond),
				dt.Retry.Round(time.Millisecond),
			)
		}(task)
	}
}

// sendMention sends a text message with @mentions.
func sendMention(ctx *Context, text string, jids []string) {
	msg := &waProto.Message{
		ExtendedTextMessage: &waProto.ExtendedTextMessage{
			Text: proto.String(text),
			ContextInfo: &waProto.ContextInfo{
				MentionedJID: jids,
			},
		},
	}
	id := ctx.Client.GenerateMessageID()
	sendQueue <- sendTask{
		client:     ctx.Client,
		to:         ctx.Event.Info.Chat,
		msg:        msg,
		id:         id,
		enqueuedAt: time.Now(),
	}
}
