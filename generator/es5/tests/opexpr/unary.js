$module('unary', function () {
  var $static = this;
  this.$class('SomeClass', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      var init = [];
      return $promise.all(init).then(function () {
        return instance;
      });
    };
    $static.$not = function (instance) {
      var $current = 0;
      var $continue = function ($resolve, $reject) {
        $resolve(instance);
        return;
      };
      return $promise.new($continue);
    };
    this.$typesig = function () {
      return $t.createtypesig(['not', 4, $g.____testlib.basictypes.Function($g.unary.SomeClass).$typeref()], ['new', 1, $g.____testlib.basictypes.Function($g.unary.SomeClass).$typeref()]);
    };
  });

  $static.DoSomething = function (first) {
    var $result;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            $g.unary.SomeClass.$not(first).then(function ($result0) {
              $result = $result0;
              $current = 1;
              $continue($resolve, $reject);
              return;
            }).catch(function (err) {
              $reject(err);
              return;
            });
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
