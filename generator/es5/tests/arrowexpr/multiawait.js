$module('multiawait', function () {
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
    $static.$plus = function (first, second) {
      var $current = 0;
      var $continue = function ($resolve, $reject) {
        $resolve(first);
        return;
      };
      return $promise.new($continue);
    };
    this.$typesig = function () {
      return $t.createtypesig(['plus', 4, $g.____testlib.basictypes.Function($g.multiawait.SomeClass).$typeref()], ['new', 1, $g.____testlib.basictypes.Function($g.multiawait.SomeClass).$typeref()]);
    };
  });

  $static.DoSomething = function (p, q) {
    var $result;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            $promise.translate(p).then(function ($result1) {
              return $promise.translate(q).then(function ($result2) {
                return $g.multiawait.SomeClass.$plus($result1, $result2).then(function ($result0) {
                  $result = $result0;
                  $current = 1;
                  $continue($resolve, $reject);
                  return;
                });
              });
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 1:
            $result;
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
