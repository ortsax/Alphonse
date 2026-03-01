package plugins

import (
	"context"
	"strings"
	"time"

	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

// Command defines a bot command.
type Command struct {
	Pattern  string
	Aliases  []string
	IsSudo   bool
	IsGroup  bool
	Category string
	Func     func(ctx *Context) error
}

// Context carries per-invocation state passed to command handlers.
type Context struct {
	Client  *whatsmeow.Client
	Event   *events.Message
	Args    []string // whitespace-split words after the command name
	Text    string   // everything after the command name (unsplit)
	Prefix  string   // the matched prefix character(s)
	Matched string   // the matched command name (lowercased)
}

// sendTimeout caps how long a single send may wait for the server ACK.
// The default (75 s) means a stuck send holds messageSendLock for 75 s,
// blocking every other Reply call.  20 s is still generous for normal
// WhatsApp RTT (~100–500 ms) but releases the lock much sooner on failures.
const sendTimeout = 20 * time.Second

// Reply sends a plain-text message back to the originating chat.
func (c *Context) Reply(text string) (whatsmeow.SendResponse, error) {
	return c.Client.SendMessage(context.Background(), c.Event.Info.Chat, &waProto.Message{
		Conversation: proto.String(text),
	}, whatsmeow.SendRequestExtra{Timeout: sendTimeout})
}

var registry []*Command

// registryMap provides O(1) command lookup by lowercased pattern or alias.
// Built incrementally during init() – read-only after startup.
var registryMap = make(map[string]*Command)

// categoryMap groups commands by their lowercased category name.
// Built incrementally during init() – read-only after startup.
var categoryMap = make(map[string][]*Command)

// Register adds a command to the global registry and updates the lookup maps.
func Register(cmd *Command) {
	registry = append(registry, cmd)
	registryMap[strings.ToLower(cmd.Pattern)] = cmd
	for _, alias := range cmd.Aliases {
		registryMap[strings.ToLower(alias)] = cmd
	}
	cat := strings.ToLower(cmd.Category)
	if cat == "" {
		cat = "general"
	}
	categoryMap[cat] = append(categoryMap[cat], cmd)
}

// parseCommand tries to extract (prefix, commandName, rest) from text using the
// configured prefixes. Matching is case-insensitive. A space between the prefix
// and the command name is allowed (e.g. ". ping" == ".ping").
func parseCommand(text string, prefixes []string) (prefix, name, rest string, ok bool) {
	lower := strings.ToLower(text)
	for _, p := range prefixes {
		var after string
		if p == "" {
			after = lower
		} else {
			lp := strings.ToLower(p)
			if !strings.HasPrefix(lower, lp) {
				continue
			}
			after = lower[len(lp):]
		}
		after = strings.TrimLeft(after, " ")
		if after == "" {
			continue
		}
		// Index-based split avoids allocating a []string.
		if i := strings.IndexByte(after, ' '); i != -1 {
			name = after[:i]
			rest = strings.TrimSpace(after[i+1:])
		} else {
			name = after
		}
		return p, name, rest, true
	}
	return "", "", "", false
}

// findCommand looks up a command by pattern or alias in O(1).
// name must already be lowercased (as returned by parseCommand).
func findCommand(name string) *Command {
	return registryMap[name]
}

// extractText returns the human-readable text from a message, checking plain
// conversation, extended text, and image/video captions.
func extractText(evt *events.Message) string {
	if t := evt.Message.GetConversation(); t != "" {
		return t
	}
	if t := evt.Message.GetExtendedTextMessage().GetText(); t != "" {
		return t
	}
	if t := evt.Message.GetImageMessage().GetCaption(); t != "" {
		return t
	}
	if t := evt.Message.GetVideoMessage().GetCaption(); t != "" {
		return t
	}
	return ""
}

// Dispatch parses a message event, resolves the matching command, enforces all
// access checks, and calls the command handler.
func Dispatch(client *whatsmeow.Client, evt *events.Message) {
	text := extractText(evt)
	if text == "" {
		return
	}

	senderID := evt.Info.Sender.User // LID user part
	isGroup := evt.Info.Chat.Server == types.GroupServer

	// Silently ignore all group messages when group chat responses are off.
	if isGroup && BotSettings.IsGCDisabled() {
		return
	}

	prefix, name, rest, ok := parseCommand(text, BotSettings.GetPrefixes())
	if !ok {
		return
	}

	cmd := findCommand(name)
	if cmd == nil {
		// If the word matches a registered category, show that category's menu.
		if menu := CategoryMenu(name); menu != "" {
			miniCtx := &Context{Client: client, Event: evt}
			miniCtx.Reply(menu)
		}
		return
	}

	ctx := &Context{
		Client:  client,
		Event:   evt,
		Args:    strings.Fields(rest),
		Text:    rest,
		Prefix:  prefix,
		Matched: name,
	}

	isSudo := BotSettings.IsSudo(senderID)
	// Also check phone number — sudoers may have been built from phone strings.
	if !isSudo && evt.Info.SenderAlt.User != "" {
		isSudo = BotSettings.IsSudo(evt.Info.SenderAlt.User)
	}

	// Banned users are silently ignored — no response given.
	isBanned := BotSettings.IsBanned(senderID)
	if !isBanned && evt.Info.SenderAlt.User != "" {
		isBanned = BotSettings.IsBanned(evt.Info.SenderAlt.User)
	}
	if isBanned {
		return
	}
	mode := BotSettings.GetMode()

	// Private mode: silently ignore non-sudo users – no response at all.
	if mode == ModePrivate && !isSudo {
		return
	}

	if cmd.IsGroup && !isGroup {
		ctx.Reply(T().GroupOnly)
		return
	}

	if cmd.IsSudo && !isSudo {
		ctx.Reply(T().SudoOnly)
		return
	}

	if BotSettings.IsCmdDisabled(name) {
		ctx.Reply(T().CmdIsDisabled)
		return
	}

	_ = cmd.Func(ctx)
}
