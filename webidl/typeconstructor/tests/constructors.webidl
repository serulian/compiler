[Constructor]
interface SomeInterface {};

[Constructor(SomeInterface value)]
interface AnotherInterface {};

[Constructor(SomeInterface value),
 Constructor(SomeInterface value, optional AnotherInterface anotherValue)]
interface ThirdInterface {};

[Constructor(SomeInterface value),
 Constructor(AnotherInterface value)]
interface FourthInterface {};