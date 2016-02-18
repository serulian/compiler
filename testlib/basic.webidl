[PrimaryGlobal]
interface Window {
	static void debugprint(any value);
}

[Constructor, NativeOperator=Plus, NativeOperator=Equals]
interface String {};

[Constructor(any value), NativeOperator=Equals]
interface Boolean {
	String toString();
};

[Constructor(optional any value),
 NativeOperator=Plus, NativeOperator=Minus, NativeOperator=Equals]
interface Number {
	String toString();
};

[Constructor]
interface Array {
	readonly attribute Number length;
	getter any (Number propertyName);
	setter void (Number propertyName, any value);
};

[Constructor]
interface Object {
  getter any (String propertyName);
  setter void (String propertyName, any value);
  static Array keys(Object o);
};