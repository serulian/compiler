[Constructor(any someparam), NativeOperator=Plus]
interface SomeBrowserThing {
	readonly attribute SomeBrowserThing InstanceAttr;
	static readonly attribute SomeBrowserThing SomeStaticAttribute;
	static SomeBrowserThing SomeStaticFunction();
	SomeBrowserThing SomeInterfaceFunction();

    getter any (any propertyName);
    setter void (any propertyName, any propertyValue);
}

[GlobalContext]
interface WindowContext {
	static readonly attribute Boolean boolValue;
}