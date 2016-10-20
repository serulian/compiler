$module('basic', function () {
  var $static = this;
  this.$class('7dbf4efd', 'SomeClass', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      var init = [];
      instance.SomeInt = $t.box(2, $g.____testlib.basictypes.Integer);
      init.push($g.basic.CoolFunction().then(function ($result0) {
        instance.AnotherBool = $result0;
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
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
        "AnotherFunction|2|29dc432d<void>": true,
      };
      return this.$cachedtypesig = computed;
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
