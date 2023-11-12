package mrcore

type (
    ClientSection interface {
        Caption() string
        Path(actionPath string) string
        MiddlewareWithPermission(name string, next HttpHandlerFunc) HttpHandlerFunc
    }

    AccessObject interface {
        CheckPrivilege(name string) bool
        CheckPermission(name string) bool
    }
)
