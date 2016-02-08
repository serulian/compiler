[Constructor]
interface String {};

[Constructor(any value), NativeOperator=Equals]
interface Boolean {};

[Constructor(optional any value),
 NativeOperator=Plus, NativeOperator=Minus, NativeOperator=Equals]
interface Number {};

[Constructor]
interface Array {
	readonly attribute Number length;
};