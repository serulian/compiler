$module('notop', function () {
  var $static = this;
  this.$class('df288f51', 'SomeClass', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function (boolValue) {
      var instance = new $static();
      instance.boolValue = boolValue;
      return $promise.resolve(instance);
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
      return {
      };
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
