$module('notop', function () {
  var $static = this;
  this.$class('SomeClass', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function (boolValue) {
      var instance = new $static();
      var init = [];
      instance.boolValue = boolValue;
      return $promise.all(init).then(function () {
        return instance;
      });
    };
    $static.$bool = function (sc) {
      var $current = 0;
      var $continue = function ($resolve, $reject) {
        $resolve(sc.boolValue);
        return;
      };
      return $promise.new($continue);
    };
    this.$typesig = function () {
      return $t.createtypesig(['bool', 4, $g.____testlib.basictypes.Function($g.____testlib.basictypes.Boolean).$typeref()], ['new', 1, $g.____testlib.basictypes.Function($g.notop.SomeClass).$typeref()]);
    };
  });

  $static.TEST = function () {
    var $result;
    var sc;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            $g.notop.SomeClass.new($t.box(false, $g.____testlib.basictypes.Boolean)).then(function ($result0) {
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
            sc = $result;
            $g.notop.SomeClass.$bool(sc).then(function ($result0) {
              $result = $t.box(!$t.unbox($result0), $g.____testlib.basictypes.Boolean);
              $current = 2;
              $continue($resolve, $reject);
              return;
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 2:
            $resolve($result);
            return;

          default:
            $resolve();
            return;
        }
      }
    };
    return $promise.new($continue);
  };
});
