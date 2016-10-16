$module('classcastfail', function () {
  var $static = this;
  this.$class('SomeClass', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      return $promise.resolve(instance);
    };
    this.$typesig = function () {
      return $t.createtypesig(['new', 1, $g.____testlib.basictypes.Function($g.classcastfail.SomeClass).$typeref()]);
    };
  });

  this.$class('AnotherClass', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      return $promise.resolve(instance);
    };
    this.$typesig = function () {
      return $t.createtypesig(['new', 1, $g.____testlib.basictypes.Function($g.classcastfail.AnotherClass).$typeref()]);
    };
  });

  $static.TEST = function () {
    var $result;
    var ac;
    var sc;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            $g.classcastfail.SomeClass.new().then(function ($result0) {
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
            $g.classcastfail.AnotherClass.new().then(function ($result0) {
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
            ac = $result;
            $t.cast(ac, $g.classcastfail.SomeClass, false);
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
});
