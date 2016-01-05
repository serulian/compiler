$module('basic', function () {
  var $static = this;
  $static.DoSomething = function () {
    var first;
    var second;
    var $state = $t.sm(function ($callback) {
      while (true) {
        switch ($state.current) {
          case 0:
            $global.SomeBrowserThing.SomeStaticAttribute;
            $global.SomeBrowserThing.SomeStaticFunction();
            $t.dynamicaccess($global.SomeBrowserThing, 'SomeStaticFunction', null, false, true);
            $global.SomeBrowserThing.SomeStaticAttribute.SomeInterfaceFunction();
            $t.dynamicaccess($global.SomeBrowserThing.SomeStaticAttribute, 'SomeInterfaceFunction', null, false, true);
            $t.nativenew($global.SomeBrowserThing, 'new', null, false, false)('foo');
            first = $t.nativenew($global.SomeBrowserThing, 'new', null, false, false)('foo');
            second = $t.nativenew($global.SomeBrowserThing, 'new', null, false, false)('bar');
            first + second;
            $state.current = -1;
            $state.returnValue = null;
            $callback($state);
            return;

          default:
            $state.current = -1;
            return;
        }
      }
    });
    return $promise.build($state);
  };
});
