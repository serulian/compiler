[Constructor]
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
};