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
  });

  $static.DoSomething = function (p, q) {
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            $promise.translate(p).then(function ($result0) {
              return $promise.translate(q).then(function ($result1) {
                return $g.multiawait.SomeClass.$plus($result0, $result1).then(function ($result2) {
                  $result = $result2;
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
