[Constructor(any someparam), NativeOperator=Plus]
interface SomeBrowserThing {
	readonly attribute SomeBrowserThing InstanceAttr;
	static readonly attribute SomeBrowserThing SomeStaticAttribute;
	static SomeBrowserThing SomeStaticFunction();
	SomeBrowserThing SomeInterfaceFunction();
}

[GlobalContext]
interface WindowContext {
	static readonly attribute any boolValue;
}