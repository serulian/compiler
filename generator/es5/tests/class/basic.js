$module('basic', function () {
  var $static = this;
  this.$class('SomeClass', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      var init = [];
      init.push(function () {
        return $promise.wrap(function () {
          $this.SomeInt = 2;
        });
      }());
      init.push(function () {
        var $returnValue$1;
        var $state = $t.sm(function ($callback) {
          while (true) {
            switch ($state.current) {
              case 0:
                $g.basic.CoolFunction().then(function (returnValue) {
                  $state.current = 1;
                  $returnValue$1 = returnValue;
                  $state.next($callback);
                }).catch(function (e) {
                  $state.error = e;
                  $state.current = -1;
                  $callback($state);
                });
                return;

              case 1:
                $this.AnotherInt = $returnValue$1;
                $state.current = -1;
                $callback($state);
                return;

              default:
                $state.current = -1;
                return;
            }
          }
        });
        return $promise.build($state);
      }());
      return $promise.all(init).then(function () {
        return instance;
      });
    };
    $instance.AnotherFunction = function () {
      return $promise.empty();
    };
  });

  $static.CoolFunction = function () {
    var $state = $t.sm(function ($callback) {
      while (true) {
        switch ($state.current) {
          case 0:
            $state.returnValue = 4;
            $state.current = -1;
            $callback($state);
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
