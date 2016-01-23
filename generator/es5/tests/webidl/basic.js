$module('basic', function () {
  var $static = this;
  $static.AnotherFunction = function () {
    return $promise.empty();
  };
  $static.DoSomething = function () {
    var first;
    var second;
    var $state = $t.sm(function ($callback) {
      while (true) {
        switch ($state.current) {
          case 0:
            $g.basic.AnotherFunction().then(function ($result0) {
              $result = $result0;
              $state.current = 1;
              $callback($state);
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 1:
            $result;
            $global.SomeBrowserThing.SomeStaticAttribute;
            $global.SomeBrowserThing.SomeStaticFunction();
            $t.dynamicaccess($global.SomeBrowserThing, 'SomeStaticFunction');
            $global.SomeBrowserThing.SomeStaticAttribute.SomeInterfaceFunction();
            $t.dynamicaccess($global.SomeBrowserThing.SomeStaticAttribute, 'SomeInterfaceFunction');
            $t.nativenew($global.SomeBrowserThing)('foo');
            $global.SomeBrowserThing.SomeStaticAttribute.SomeInterfaceFunction().InstanceAttr;
            first = $t.nativenew($global.SomeBrowserThing)('foo');
            second = $t.nativenew($global.SomeBrowserThing)('bar');
            first + second;
            $state.current = -1;
            return;

          default:
            $state.current = -1;
            return;
        }
      }
    });
    return $promise.build($state);
  };
  $static.TEST = function () {
    var $state = $t.sm(function ($callback) {
      while (true) {
        switch ($state.current) {
          case 0:
            $state.resolve($global.boolValue);
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
