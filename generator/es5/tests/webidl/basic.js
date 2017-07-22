$module('basic', function () {
  var $static = this;
  $static.AnotherFunction = function () {
    return;
  };
  $static.DoSomething = function () {
    var first;
    var second;
    $g.basic.AnotherFunction();
    $global.SomeBrowserThing.SomeStaticAttribute;
    $global.SomeBrowserThing.SomeStaticFunction();
    $t.dynamicaccess($global.SomeBrowserThing, 'SomeStaticFunction', false);
    $global.SomeBrowserThing.SomeStaticAttribute.SomeInterfaceFunction();
    $t.dynamicaccess($global.SomeBrowserThing.SomeStaticAttribute, 'SomeInterfaceFunction', false);
    $t.nativenew($global.SomeBrowserThing)($t.fastbox('foo', $g.________testlib.basictypes.String));
    $global.SomeBrowserThing.SomeStaticAttribute.SomeInterfaceFunction().InstanceAttr;
    first = $t.nativenew($global.SomeBrowserThing)($t.fastbox('foo', $g.________testlib.basictypes.String));
    second = $t.nativenew($global.SomeBrowserThing)($t.fastbox('bar', $g.________testlib.basictypes.String));
    first[$t.fastbox('hello', $g.________testlib.basictypes.String)] = second;
    return;
  };
  $static.TEST = function () {
    return $t.fastbox($global.boolValue, $g.________testlib.basictypes.Boolean);
  };
});
