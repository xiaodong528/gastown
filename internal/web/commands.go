package web

import (
	"fmt"
	"regexp"
	"strings"
)

// CommandMeta describes a command's properties for the dashboard.
type CommandMeta struct {
	// Safe commands can run without user confirmation
	Safe bool
	// Confirm commands require user confirmation before execution
	Confirm bool
	// Desc is a short description shown in the command palette
	Desc string
	// Category groups commands in the palette UI
	Category string
	// Args is a placeholder hint for required arguments (e.g., "<convoy-id>")
	Args string
	// ArgType specifies what kind of options to show (rigs, polecats, convoys, agents, hooks)
	ArgType string
}

// AllowedCommands defines which gt commands can be executed from the dashboard.
// Commands not in this list are blocked for security.
var AllowedCommands = map[string]CommandMeta{
	// === Read-only commands (always safe) ===
	"status":      {Safe: true, Desc: "Show workspace status（显示工作空间状态）", Category: "Status（状态）"},
	"agents list": {Safe: true, Desc: "List active agents（列出活跃智能体）", Category: "Status（状态）"},
	"convoy list": {Safe: true, Desc: "List Convoy（列出 Convoy）", Category: "Convoy"},
	"convoy show":   {Safe: true, Desc: "Show Convoy details（查看 Convoy 详情）", Category: "Convoy", Args: "<convoy-id>", ArgType: "convoys"},
	"convoy status": {Safe: true, Desc: "Show Convoy status and tracked issues（显示 Convoy 状态与跟踪的 Issue）", Category: "Convoy", Args: "<convoy-id> --json", ArgType: "convoys"},
	"mail inbox":  {Safe: true, Desc: "View inbox（查看收件箱）", Category: "Mail（邮件）"},
	"mail check":  {Safe: true, Desc: "Check new mail（检查新邮件）", Category: "Mail（邮件）"},
	"mail peek":   {Safe: true, Desc: "Preview message（预览消息）", Category: "Mail（邮件）", Args: "<message-id>"},
	"rig list":    {Safe: true, Desc: "List Rig（列出 Rig）", Category: "Rig"},
	"rig show":    {Safe: true, Desc: "Show Rig details（查看 Rig 详情）", Category: "Rig", Args: "<rig-name>", ArgType: "rigs"},
	"doctor":      {Safe: true, Desc: "Health check（健康检查）", Category: "Diagnostics（诊断）"},
	"hooks list":  {Safe: true, Desc: "List Hook（列出 Hook）", Category: "Hook"},
	"activity":    {Safe: true, Desc: "Show recent activity（显示近期活动）", Category: "Status（状态）"},
	"info":        {Safe: true, Desc: "Show workspace info（显示工作空间信息）", Category: "Status（状态）"},
	"log":         {Safe: true, Desc: "View logs（查看日志）", Category: "Diagnostics（诊断）"},
	"audit":       {Safe: true, Desc: "View audit logs（查看审计日志）", Category: "Diagnostics（诊断）"},

	// Polecat read-only
	"polecat list --all": {Safe: true, Desc: "List all Polecat（列出所有 Polecat）", Category: "Polecat"},
	"polecat show":       {Safe: true, Desc: "Show Polecat details（查看 Polecat 详情）", Category: "Polecat", Args: "<rig>/<name>", ArgType: "polecats"},

	// Crew read-only
	"crew list --all": {Safe: true, Desc: "List all Crew members（列出所有 Crew 成员）", Category: "Crew"},
	"crew show":       {Safe: true, Desc: "Show Crew details（查看 Crew 详情）", Category: "Crew", Args: "<rig>/<name>", ArgType: "crew"},

	// === Action commands (require confirmation) ===

	// Mail actions
	"mail send":      {Confirm: true, Desc: "Send message（发送消息）", Category: "Mail（邮件）", Args: "<address> -s <subject> -m <message>", ArgType: "agents"},
	"mail mark-read": {Confirm: false, Desc: "Mark as read（标记为已读）", Category: "Mail（邮件）", Args: "<message-id>", ArgType: "messages"},
	"mail archive":   {Confirm: false, Desc: "Archive message（归档消息）", Category: "Mail（邮件）", Args: "<message-id>", ArgType: "messages"},
	"mail reply":     {Confirm: true, Desc: "Reply to message（回复消息）", Category: "Mail（邮件）", Args: "<message-id> -m <message>", ArgType: "messages"},

	// Escalation actions
	"escalate ack":      {Confirm: true, Desc: "Acknowledge escalation（确认升级事件）", Category: "Escalation（升级事件）", Args: "<escalation-id>", ArgType: "escalations"},
	"escalate resolve":  {Confirm: true, Desc: "Resolve escalation（解决升级事件）", Category: "Escalation（升级事件）", Args: "<escalation-id>", ArgType: "escalations"},
	"escalate reassign": {Confirm: true, Desc: "Reassign escalation（重新分配升级事件）", Category: "Escalation（升级事件）", Args: "<escalation-id> <agent>", ArgType: "escalations"},

	// Convoy actions
	"convoy create":  {Confirm: true, Desc: "Create Convoy（创建 Convoy）", Category: "Convoy", Args: "<name>"},
	"convoy refresh": {Confirm: false, Desc: "Refresh Convoy（刷新 Convoy）", Category: "Convoy", Args: "<convoy-id>", ArgType: "convoys"},
	"convoy add":     {Confirm: true, Desc: "Add issue to Convoy（添加 Issue 到 Convoy）", Category: "Convoy", Args: "<convoy-id> <issue>", ArgType: "convoys"},

	// Rig actions
	"rig boot":  {Confirm: true, Desc: "Boot Rig（启动 Rig）", Category: "Rig", Args: "<rig-name>", ArgType: "rigs"},
	"rig start": {Confirm: true, Desc: "Start Rig（启动 Rig）", Category: "Rig", Args: "<rig-name>", ArgType: "rigs"},

	// Agent lifecycle (careful)
	"witness start":  {Confirm: true, Desc: "Start Witness（启动 Witness）", Category: "Agent（智能体）", Args: "<rig-name>", ArgType: "rigs"},
	"refinery start": {Confirm: true, Desc: "Start Refinery（启动 Refinery）", Category: "Agent（智能体）", Args: "<rig-name>", ArgType: "rigs"},
	"mayor attach":   {Confirm: true, Desc: "Attach Mayor（连接 Mayor）", Category: "Agent（智能体）"},
	"deacon start":   {Confirm: true, Desc: "Start Deacon（启动 Deacon）", Category: "Agent（智能体）"},

	// Polecat actions
	"polecat add":    {Confirm: true, Desc: "Add Polecat（添加 Polecat）", Category: "Polecat", Args: "<rig> <name>", ArgType: "rigs"},
	"polecat remove": {Confirm: true, Desc: "Remove Polecat（移除 Polecat）", Category: "Polecat", Args: "<rig>/<name>", ArgType: "polecats"},

	// Work assignment
	"sling":       {Confirm: true, Desc: "Assign work to agent（分配工作给智能体）", Category: "Work（工作）", Args: "<bead> <rig>", ArgType: "hooks"},
	"unsling":     {Confirm: true, Desc: "Unassign work（取消分配工作）", Category: "Work（工作）", Args: "<bead>", ArgType: "hooks"},
	"hook attach": {Confirm: true, Desc: "Mount Hook（挂载 Hook）", Category: "Hook", Args: "<bead>", ArgType: "hooks"},
	"hook detach": {Confirm: true, Desc: "Detach Hook（卸载 Hook）", Category: "Hook", Args: "<bead>", ArgType: "hooks"},

	// Notifications
	"notify":    {Confirm: true, Desc: "Send notification（发送通知）", Category: "Notification（通知）", Args: "<message>"},
	"broadcast": {Confirm: true, Desc: "Broadcast message（广播消息）", Category: "Notification（通知）", Args: "<message>"},
}

// BlockedPatterns are regex patterns for commands that should never run from the dashboard.
// These require terminal access for safety.
var BlockedPatterns = []*regexp.Regexp{
	regexp.MustCompile(`--force`),
	regexp.MustCompile(`--hard`),
	regexp.MustCompile(`\brm\b`),
	regexp.MustCompile(`\bdelete\b`),
	regexp.MustCompile(`\bkill\b`),
	regexp.MustCompile(`\bdestroy\b`),
	regexp.MustCompile(`\bpurge\b`),
	regexp.MustCompile(`\breset\b`),
	regexp.MustCompile(`\bclean\b`),
}

// ValidateCommand checks if a command is allowed to run from the dashboard.
// Returns the command metadata if allowed, or an error if blocked.
func ValidateCommand(rawCommand string) (*CommandMeta, error) {
	rawCommand = strings.TrimSpace(rawCommand)
	if rawCommand == "" {
		return nil, fmt.Errorf("empty command")
	}

	// Check blocked patterns first
	for _, pattern := range BlockedPatterns {
		if pattern.MatchString(rawCommand) {
			return nil, fmt.Errorf("command contains blocked pattern: %s", pattern.String())
		}
	}

	// Extract base command (first 1-2 words) for whitelist lookup
	baseCmd := extractBaseCommand(rawCommand)

	meta, ok := AllowedCommands[baseCmd]
	if !ok {
		return nil, fmt.Errorf("command not in whitelist: %s", baseCmd)
	}

	return &meta, nil
}

// extractBaseCommand gets the command prefix for whitelist matching.
// "mail send foo bar" -> "mail send"
// "status --json" -> "status"
// "polecat list --all" -> "polecat list --all"
func extractBaseCommand(cmd string) string {
	parts := strings.Fields(cmd)
	if len(parts) == 0 {
		return ""
	}

	// Try three-word command first (e.g., "polecat list --all")
	if len(parts) >= 3 {
		threeWord := parts[0] + " " + parts[1] + " " + parts[2]
		if _, ok := AllowedCommands[threeWord]; ok {
			return threeWord
		}
	}

	// Try two-word command (e.g., "convoy list")
	if len(parts) >= 2 {
		twoWord := parts[0] + " " + parts[1]
		if _, ok := AllowedCommands[twoWord]; ok {
			return twoWord
		}
	}

	// Fall back to single-word command
	return parts[0]
}

// SanitizeArgs removes potentially dangerous characters from command arguments.
// This is a defense-in-depth measure; the whitelist is the primary protection.
func SanitizeArgs(args []string) []string {
	sanitized := make([]string, 0, len(args))
	for _, arg := range args {
		// Remove shell metacharacters
		clean := strings.Map(func(r rune) rune {
			switch r {
			case ';', '|', '&', '$', '`', '(', ')', '{', '}', '<', '>', '\n', '\r':
				return -1 // Remove character
			default:
				return r
			}
		}, arg)
		if clean != "" {
			sanitized = append(sanitized, clean)
		}
	}
	return sanitized
}

// GetCommandList returns all allowed commands for the command palette UI.
func GetCommandList() []CommandInfo {
	commands := make([]CommandInfo, 0, len(AllowedCommands))
	for name, meta := range AllowedCommands {
		commands = append(commands, CommandInfo{
			Name:     name,
			Desc:     meta.Desc,
			Category: meta.Category,
			Safe:     meta.Safe,
			Confirm:  meta.Confirm,
			Args:     meta.Args,
			ArgType:  meta.ArgType,
		})
	}
	return commands
}

// CommandInfo is the JSON-serializable form of a command for the UI.
type CommandInfo struct {
	Name     string `json:"name"`
	Desc     string `json:"desc"`
	Category string `json:"category"`
	Safe     bool   `json:"safe"`
	Confirm  bool   `json:"confirm"`
	Args     string `json:"args,omitempty"`
	ArgType  string `json:"argType,omitempty"`
}
