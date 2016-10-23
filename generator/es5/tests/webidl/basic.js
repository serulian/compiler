$module('basic', function () {
  var $static = this;
  $static.AnotherFunction = function () {
    return $promise.empty();
  };
  $static.DoSomething = function () {
    var $result;
    var first;
    var second;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            $g.basic.AnotherFunction().then(function ($result0) {
              $result = $result0;
              $current = 1;
              $continue($resolve, $reject);
              return;
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 1:
            $global.SomeBrowserThing.SomeStaticAttribute;
            $global.SomeBrowserThing.SomeStaticFunction();
            $t.dynamicaccess($global.SomeBrowserThing, 'SomeStaticFunction').then(function ($result0) {
              $result = $result0;
              $current = 2;
              $continue($resolve, $reject);
              return;
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 2:
            $global.SomeBrowserThing.SomeStaticAttribute.SomeInterfaceFunction();
            $t.dynamicaccess($global.SomeBrowserThing.SomeStaticAttribute, 'SomeInterfaceFunction').then(function ($result0) {
              $result = $result0;
              $current = 3;
              $continue($resolve, $reject);
              return;
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 3:
            $t.nativenew($global.SomeBrowserThing)($t.fastbox('foo', $g.____testlib.basictypes.String));
            $global.SomeBrowserThing.SomeStaticAttribute.SomeInterfaceFunction().InstanceAttr;
            first = $t.nativenew($global.SomeBrowserThing)($t.fastbox('foo', $g.____testlib.basictypes.String));
            second = $t.nativenew($global.SomeBrowserThing)($t.fastbox('bar', $g.____testlib.basictypes.String));
            first[$t.fastbox('hello', $g.____testlib.basictypes.String)] = second;
            $resolve();
            return;

          default:
            $resolve();
            return;
        }
      }
    };
    return $promise.new($continue);
  };
  $static.TEST = function () {
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      $resolve($t.fastbox($global.boolValue, $g.____testlib.basictypes.Boolean));
      return;
    };
    return $promise.new($continue);
  };
});
