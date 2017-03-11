interface ISomeCollapsedType {
	readonly attribute any Third;

	// repeat First with the same signature, which should be ignored.
	readonly attribute any First;	
};