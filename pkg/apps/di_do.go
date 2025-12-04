package apps

import (
	"github.com/samber/do/v2"
)

type Injector do.Injector

var container = do.New()

func appContainer() Injector {
	return container
}

type Provider[T any] do.Provider[T]

func NameOf[T any]() string {
	return do.NameOf[T]()
}

func Provide[T any](provider Provider[T]) {
	do.Provide[T](appContainer(), do.Provider[T](provider))
}

func ProvideNamed[T any](name string, provider Provider[T]) {
	do.ProvideNamed[T](appContainer(), name, do.Provider[T](provider))
}

func ProvideValue[T any](value T) {
	do.ProvideValue[T](appContainer(), value)
}

func ProvideNamedValue[T any](name string, value T) {
	do.ProvideNamedValue[T](appContainer(), name, value)
}

func ProvideTransient[T any](provider Provider[T]) {
	do.ProvideTransient[T](appContainer(), do.Provider[T](provider))
}

func ProvideNamedTransient[T any](name string, provider Provider[T]) {
	do.ProvideNamedTransient[T](appContainer(), name, do.Provider[T](provider))
}

func Override[T any](provider Provider[T]) {
	do.Override[T](appContainer(), do.Provider[T](provider))
}

func OverrideNamed[T any](name string, provider Provider[T]) {
	do.OverrideNamed[T](appContainer(), name, do.Provider[T](provider))
}

func OverrideValue[T any](value T) {
	do.OverrideValue[T](appContainer(), value)
}

func OverrideNamedValue[T any](name string, value T) {
	do.OverrideNamedValue[T](appContainer(), name, value)
}

func OverrideTransient[T any](provider Provider[T]) {
	do.OverrideTransient[T](appContainer(), do.Provider[T](provider))
}

func OverrideNamedTransient[T any](name string, provider Provider[T]) {
	do.OverrideNamedTransient[T](appContainer(), name, do.Provider[T](provider))
}

func Invoke[T any]() (T, error) {
	return do.Invoke[T](appContainer())
}

func InvokeNamed[T any](name string) (T, error) {
	return do.InvokeNamed[T](appContainer(), name)
}

func MustInvoke[T any]() T {
	return do.MustInvoke[T](appContainer())
}

func MustInvokeNamed[T any](name string) T {
	return do.MustInvokeNamed[T](appContainer(), name)
}

func InvokeStruct[T any]() (T, error) {
	return do.InvokeStruct[T](appContainer())
}

func MustInvokeStruct[T any]() T {
	return do.MustInvokeStruct[T](appContainer())
}

func As[Initial any, Alias any]() error {
	return do.As[Initial, Alias](appContainer())
}

func MustAs[Initial any, Alias any]() {
	do.MustAs[Initial, Alias](appContainer())
}

func AsNamed[Initial any, Alias any](initial string, alias string) error {
	return do.AsNamed[Initial, Alias](appContainer(), initial, alias)
}

func MustAsNamed[Initial any, Alias any](initial string, alias string) {
	do.MustAsNamed[Initial, Alias](appContainer(), initial, alias)
}

func InvokeAs[T any]() (T, error) {
	return do.InvokeAs[T](appContainer())
}

func MustInvokeAs[T any]() T {
	return do.MustInvokeAs[T](appContainer())
}

func Package(services ...func(Injector)) func(Injector) {
	return func(injector Injector) {
		for i := range services {
			services[i](injector)
		}
	}
}

func Lazy[T any](p Provider[T]) func(Injector) {
	return func(injector Injector) {
		do.Provide[T](injector, do.Provider[T](p))
	}
}

func LazyNamed[T any](serviceName string, p Provider[T]) func(Injector) {
	return func(injector Injector) {
		do.ProvideNamed[T](injector, serviceName, do.Provider[T](p))
	}
}

func Eager[T any](value T) func(Injector) {
	return func(injector Injector) {
		do.ProvideValue[T](injector, value)
	}
}

func EagerNamed[T any](serviceName string, value T) func(Injector) {
	return func(injector Injector) {
		do.ProvideNamedValue[T](injector, serviceName, value)
	}
}

func Transient[T any](p Provider[T]) func(Injector) {
	return func(injector Injector) {
		do.ProvideTransient[T](injector, do.Provider[T](p))
	}
}

func TransientNamed[T any](serviceName string, p Provider[T]) func(Injector) {
	return func(injector Injector) {
		do.ProvideNamedTransient[T](injector, serviceName, do.Provider[T](p))
	}
}

func Bind[Initial any, Alias any]() func(Injector) {
	return func(injector Injector) {
		do.MustAs[Initial, Alias](injector)
	}
}

func BindNamed[Initial any, Alias any](initial string, alias string) func(Injector) {
	return func(injector Injector) {
		do.MustAsNamed[Initial, Alias](injector, initial, alias)
	}
}
