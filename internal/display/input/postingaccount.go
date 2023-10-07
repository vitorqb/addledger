package input

import (
	"github.com/vitorqb/addledger/internal/controller"
	"github.com/vitorqb/addledger/internal/display/widgets"
	eventbusmod "github.com/vitorqb/addledger/internal/eventbus"
)

func NewPostingAccount(
	controller controller.IInputController,
	eventbus eventbusmod.IEventBus,
) *widgets.InputField {
	field := widgets.NewInputField()
	field.SetLabel("Account: ")
	field.SetChangedFunc(controller.OnPostingAccountChanged)
	field.LinkContextualList(eventbus, widgets.ContextualListLinkOpts{
		InputName:           "postingaccount",
		OnListAction:        controller.OnPostingAccountListAcction,
		OnDone:              controller.OnPostingAccountDone,
		OnInsertFromContext: controller.OnPostingAccountInsertFromContext,
	})
	return field
}
