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
      var $state = $t.sm(function ($continue) {
        $state.resolve(first);
        return;
      });
      return $promise.build($state);
    };
  });

  $static.DoSomething = function (p, q) {
    var $state = $t.sm(function ($continue) {
      while (true) {
        switch ($state.current) {
          case 0:
            $promise.translate(p).then(function ($result0) {
              return $promise.translate(q).then(function ($result1) {
                return $g.multiawait.SomeClass.$plus($result0, $result1).then(function ($result2) {
                  $result = $result2;
                  $state.current = 1;
                  $continue($state);
                });
              });
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 1:
            $result;
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
});
