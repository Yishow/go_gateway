package wirebudget

import "context"

type admissionKey struct{}
type sendObserverKey struct{}

// WithAdmission binds a short request-eligibility check to every actual send.
// Checks compose parent-first and must not perform I/O or wait for a device gate.
// A successful check is the admission linearization point: later invalidation
// cannot retract that send, but every subsequent send must check again.
func WithAdmission(ctx context.Context, check func() error) context.Context {
	if check == nil {
		return ctx
	}
	previous, hasPrevious := ctx.Value(admissionKey{}).(func() error)
	return context.WithValue(ctx, admissionKey{}, func() error {
		if hasPrevious && previous != nil {
			if err := previous(); err != nil {
				return err
			}
		}
		return check()
	})
}

// WithSendObserver observes a successfully admitted, budgeted send immediately
// before transport I/O. It must be short, nonblocking and read only atomic state;
// no provider, scheduler or connection-lock calls are allowed. Observers compose
// parent-first. A diagnostic panic cannot change the acquired send allowance.
func WithSendObserver(ctx context.Context, observe func()) context.Context {
	if observe == nil {
		return ctx
	}
	previous, hasPrevious := ctx.Value(sendObserverKey{}).(func())
	return context.WithValue(ctx, sendObserverKey{}, func() {
		if hasPrevious && previous != nil {
			previous()
		}
		observeSendSafely(observe)
	})
}

func checkAdmission(ctx context.Context) error {
	if check, ok := ctx.Value(admissionKey{}).(func() error); ok && check != nil {
		return check()
	}
	return nil
}

func observeSend(ctx context.Context) {
	if observe, ok := ctx.Value(sendObserverKey{}).(func()); ok && observe != nil {
		observe()
	}
}

func observeSendSafely(observe func()) {
	defer func() { recover() }() //nolint:errcheck // An observer cannot alter acquisition.
	observe()
}
