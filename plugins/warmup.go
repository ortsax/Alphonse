package plugins

import (
	"context"
	"time"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types"
)

// StartWarmup kicks off a background goroutine that pre-establishes Signal
// sessions for all known contacts.  Call this once after the client connects
// and is fully logged in.  It is safe to run while the bot is already handling
// messages.
func StartWarmup(client *whatsmeow.Client) {
	go runWarmup(client)
}

// warmupBatchSize is the number of users submitted to GetUserDevices per usync
// request.  WhatsApp's usync endpoint accepts up to ~50 JIDs per call.
const warmupBatchSize = 50

// warmupBatchDelay is the pause between batches to avoid usync rate-limiting.
const warmupBatchDelay = 150 * time.Millisecond

func runWarmup(client *whatsmeow.Client) {
	ctx := context.Background()

	contacts, err := client.Store.Contacts.GetAllContacts(ctx)
	if err != nil {
		return
	}

	// Collect individual user JIDs (groups use sender keys, not prekeys).
	var users []types.JID
	for jid := range contacts {
		if jid.Server == types.DefaultUserServer && jid.Device == 0 {
			users = append(users, jid)
		}
	}
	if len(users) == 0 {
		return
	}
	warmed, batches := 0, 0

	for i := 0; i < len(users); i += warmupBatchSize {
		end := min(i + warmupBatchSize, len(users))
		batch := users[i:end]
		batches++

		if err := client.WarmSessions(ctx, batch); err != nil {
		} else {
			warmed += len(batch)
		}

		if end < len(users) {
			time.Sleep(warmupBatchDelay)
		}
	}
}

// warmupSender warms the Signal session for a single user immediately when a
// message is received from them.  Call in a goroutine from the event handler.
// This ensures that by the time any reply is queued, the session is already
// cached in SQLite for subsequent sends.
func warmupSender(client *whatsmeow.Client, sender types.JID) {
	// Only warm user JIDs (not groups — group encryption uses sender keys).
	if sender.Server != types.DefaultUserServer {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = client.WarmSessions(ctx, []types.JID{sender.ToNonAD()})
}
