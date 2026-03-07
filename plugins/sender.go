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
	client *whatsmeow.Client
	to     types.JID
	msg    *waProto.Message
	id     types.MessageID
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
			_, err := t.client.SendMessage(
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
		client: ctx.Client,
		to:     ctx.Event.Info.Chat,
		msg:    msg,
		id:     id,
	}
}
