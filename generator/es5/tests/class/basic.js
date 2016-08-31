$module('basic', function () {
  var $static = this;
  this.$class('SomeClass', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      var init = [];
      init.push($promise.resolve($t.box(2, $g.____testlib.basictypes.Integer)).then(function (result) {
        instance.SomeInt = result;
      }));
      init.push($g.basic.CoolFunction().then(function ($result0) {
        return $promise.resolve($result0);
      }).then(function (result) {
        instance.AnotherBool = result;
      }));
      return $promise.all(init).then(function () {
        return instance;
      });
    };
    $instance.AnotherFunction = function () {
      var $this = this;
      return $promise.empty();
    };
    this.$typesig = function () {
      return $t.createtypesig(['AnotherFunction', 2, $g.____testlib.basictypes.Function($t.void).$typeref()], ['new', 1, $g.____testlib.basictypes.Function($g.basic.SomeClass).$typeref()]);
    };
  });

  $static.CoolFunction = function () {
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      $resolve($t.box(true, $g.____testlib.basictypes.Boolean));
      return;
    };
    return $promise.new($continue);
  };
  $static.TEST = function () {
    var $result;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            $g.basic.SomeClass.new().then(function ($result0) {
              $result = $result0.AnotherBool;
              $current = 1;
              $continue($resolve, $reject);
              return;
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 1:
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
