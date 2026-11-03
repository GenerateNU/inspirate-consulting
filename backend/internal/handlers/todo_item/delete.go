package todoitem

import "context"

func (h *Handler) DeleteTodoItem(ctx context.Context, id string) error {
	err := h.TodoItemRepository.DeleteTodoItem(ctx, id)
	if(err != nil) {
		return err
	}
	return nil
}
