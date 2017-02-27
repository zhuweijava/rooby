package initializer

import (
	"fmt"
	"github.com/st0012/rooby/ast"
	"github.com/st0012/rooby/object"
)

var (
	ObjectClass = initializeObjectClass()
	ClassClass  = initializeClassClass()
	NullClass   = initializeNullClass()
)

var BuiltinGlobalMethods = []*object.BuiltInMethod{
	{
		Fn: func(receiver object.Object) object.BuiltinMethodBody {
			return func(args ...object.Object) object.Object {
				for _, arg := range args {
					fmt.Println(arg.Inspect())
				}

				return NULL
			}
		},
		Name: "puts",
	},
	{
		Fn: func(receiver object.Object) object.BuiltinMethodBody {
			return func(args ...object.Object) object.Object {
				switch r := receiver.(type) {
				case object.BaseObject:
					return r.ReturnClass()
				case object.Class:
					return r.ReturnClass()
				default:
					return &object.Error{Message: fmt.Sprint("Can't call class on %T", r)}
				}
			}
		},
		Name: "class",
	},
}

var BuiltinClassMethods = []*object.BuiltInMethod{
	{
		Fn: func(receiver object.Object) object.BuiltinMethodBody {
			return func(args ...object.Object) object.Object {
				class := receiver.(*object.RClass)
				instance := InitializeInstance(class)
				initMethod := class.LookupInstanceMethod("initialize")

				if initMethod != nil {
					instance.InitializeMethod = initMethod.(*object.Method)
				}

				return instance
			}
		},
		Name: "new",
	},
	{
		Fn: func(receiver object.Object) object.BuiltinMethodBody {
			return func(args ...object.Object) object.Object {
				name := receiver.(object.Class).ReturnName()
				nameString := &object.StringObject{Value: name.Value}
				return nameString
			}
		},
		Name: "name",
	},
}

func InitializeProgram() {
	initializeStringClass()
	initializeIntegerClass()
	initializeBooleanClass()
}

func InitializeMainObject() *object.RObject {
	obj := &object.RObject{Class: ObjectClass, InstanceVariables: object.NewEnvironment()}
	scope := &object.Scope{Self: obj, Env: object.NewEnvironment()}
	obj.Scope = scope
	return obj
}

func initializeObjectClass() *object.RClass {
	name := &ast.Constant{Value: "Object"}
	globalMethods := object.NewEnvironment()

	for _, m := range BuiltinGlobalMethods {
		globalMethods.Set(m.Name, m)
	}

	class := &object.RClass{BaseClass: &object.BaseClass{Name: name, Class: ClassClass, Methods: globalMethods}}
	NullClass.SuperClass = class
	return class
}

func initializeClassClass() *object.RClass {
	methods := object.NewEnvironment()

	for _, m := range BuiltinClassMethods {
		methods.Set(m.Name, m)
	}

	name := &ast.Constant{Value: "Class"}
	class := &object.RClass{BaseClass: &object.BaseClass{Name: name, Methods: methods}}

	return class
}

func InitializeClass(name *ast.Constant, scope *object.Scope) *object.RClass {
	class := &object.RClass{BaseClass: &object.BaseClass{Name: name, Methods: object.NewEnvironment(), Class: ClassClass, SuperClass: ObjectClass}}
	classScope := &object.Scope{Self: class, Env: object.NewClosedEnvironment(scope.Env)}
	class.Scope = classScope

	return class
}

func InitializeInstance(c *object.RClass) *object.RObject {
	instance := &object.RObject{Class: c, InstanceVariables: object.NewEnvironment()}

	return instance
}

func initializeNullClass() *object.NullClass {
	n := &ast.Constant{Value: "Null"}
	baseClass := &object.BaseClass{Name: n, Methods: object.NewEnvironment(), Class: ClassClass}
	nc := &object.NullClass{BaseClass: baseClass}
	return nc
}

func checkArgumentLen(args []object.Object, class object.Class, method_name string) *object.Error {
	if len(args) > 1 {
		return &object.Error{Message: fmt.Sprintf("Too many arguments for %s#%s", class.ReturnName().Value, method_name)}
	}

	return nil
}

func wrongTypeError(c object.Class) *object.Error {
	return &object.Error{Message: fmt.Sprintf("expect argument to be %s type", c.ReturnName().Value)}
}