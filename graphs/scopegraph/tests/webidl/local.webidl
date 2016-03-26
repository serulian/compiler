[Constructor(any someParam),
 Constructor(any someParam, any anotherParam),
 NativeOperator=Plus]
interface SomeType {
	any SomeFunction();
	static any StaticFunction();

    getter Second (First propertyName);
    setter void (Third propertyName, Fourth propertyValue);
};

interface First {};
interface Second {};
interface Third {};
interface Fourth {};

interface AnotherFirst : First {};