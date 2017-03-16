[PrimaryGlobal]
interface Window {
	static void debugprint(any value);
};

[NoInterfaceObject]
interface __serulian_internal {
  static any autoBox(any value);
  static any autoUnbox();
};

[Constructor, NativeOperator=Plus, NativeOperator=Equals]
interface String {
	serializer;
	readonly attribute Number length;
};

[Constructor(any value), NativeOperator=Equals]
interface Boolean {
	String toString();
	serializer;
};

[Constructor(optional any value),
 NativeOperator=Plus, NativeOperator=Minus, NativeOperator=Equals]
interface Number {
	String toString();
	serializer;
};

[Constructor]
interface Array {
	readonly attribute Number length;
	getter any (Number propertyName);
	setter void (Number propertyName, any value);
	Array slice(any start, any end);
	serializer;
};

[Constructor]
interface Object {
  getter any (String propertyName);
  setter void (String propertyName, any value);
  static Array keys(Object o);
  serializer;
};

[NoInterfaceObject]
interface JSON {
  static Object parse(any text, optional any reviver);
  static String stringify(any value, optional any replacer, optional any space);
};


[Constructor(optional String message, optional String fileName, optional Number lineNumber)]
interface Error {
	readonly attribute String message;
	readonly attribute String fileName;
	readonly attribute String name;
	readonly attribute Number lineNumber;
	readonly attribute Number columnNumber;
	readonly attribute String stack;
};