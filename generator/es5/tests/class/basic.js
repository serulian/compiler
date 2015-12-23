$module('basic', function () {
  var $static = this;
  this.cls('SomeClass', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function ($callback) {
      var instance = new $static();
      var init = [];
      init.push(function () {
        return $promise.wrap(function () {
          $this.SomeInt = 2;
        });
      }());
      init.push(function () {
        var $state = {
          current: 0,
          returnValue: null,
        };
        var $returnValue$1;
        $state.next = function ($callback) {
          try {
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
                  $state.returnValue = $returnValue$1;
                  $state.current = -1;
                  $callback($state);
                  return;

                default:
                  $state.current = -1;
                  return;
              }
            }
          } catch (e) {
            $state.error = e;
            $state.current = -1;
            $callback($state);
          }
        };
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
    var $state = {
      current: 0,
      returnValue: null,
    };
    $state.next = function ($callback) {
      try {
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
      } catch (e) {
        $state.error = e;
        $state.current = -1;
        $callback($state);
      }
    };
    return $promise.build($state);
  };
});
