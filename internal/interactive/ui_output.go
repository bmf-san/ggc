package interactive

import (
	"fmt"
	"time"
)

// writeError writes an error message to stderr
func (ui *UI) writeError(format string, a ...interface{}) {
	_, _ = fmt.Fprintf(ui.stderr, format+"\n", a...)
}

// write writes a message to stdout
func (ui *UI) write(format string, a ...interface{}) {
	_, _ = fmt.Fprintf(ui.stdout, format, a...)
}

// writeColor writes a colored message to stdout
func (ui *UI) writeColor(text string) {
	_, _ = fmt.Fprint(ui.stdout, text)
}

// writeln writes a message with newline to stdout
func (ui *UI) writeln(format string, a ...interface{}) {
	// Move to line start, clear line, write content, then CRLF
	_, _ = fmt.Fprint(ui.stdout, "\r\x1b[K")
	_, _ = fmt.Fprintf(ui.stdout, format+"\r\n", a...)
}

// notifySoftCancel sets the soft cancel flash notification
func (ui *UI) notifySoftCancel() {
	ui.softCancelFlash.Store(true)
}

// consumeSoftCancelFlash checks and clears the soft cancel flash notification
func (ui *UI) consumeSoftCancelFlash() bool {
	return ui.softCancelFlash.Swap(false)
}

// notifyWorkflowError displays an error message for the specified duration
func (ui *UI) notifyWorkflowError(message string, duration time.Duration) {
	if ui == nil {
		return
	}
	ui.workflowNotice = ""
	ui.workflowError = message
	ui.errorExpiresAt = time.Now().Add(duration)
}

// workflowErrorMessage returns the current error message if not expired
func (ui *UI) workflowErrorMessage() string {
	if ui == nil || ui.workflowError == "" {
		return ""
	}
	if time.Now().After(ui.errorExpiresAt) {
		ui.workflowError = ""
		return ""
	}
	return ui.workflowError
}

// notifyWorkflowSuccess displays a success message for the specified duration
func (ui *UI) notifyWorkflowSuccess(message string, duration time.Duration) {
	if ui == nil {
		return
	}
	ui.workflowError = ""
	ui.workflowNotice = message
	ui.noticeExpiresAt = time.Now().Add(duration)
}

// workflowNoticeMessage returns the current notice message if not expired
func (ui *UI) workflowNoticeMessage() string {
	if ui == nil || ui.workflowNotice == "" {
		return ""
	}
	if time.Now().After(ui.noticeExpiresAt) {
		ui.workflowNotice = ""
		return ""
	}
	return ui.workflowNotice
}
