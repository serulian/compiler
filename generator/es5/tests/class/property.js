$module('property', function () {
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
      return $promise.all(init).then(function () {
        return instance;
      });
    };
    $instance.SomeProp = function (opt_val) {
      if (arguments.length == 0) {
        return function () {
          var $this = this;
          var $state = {
            current: 0,
            returnValue: null,
          };
          $state.next = function ($callback) {
            try {
              while (true) {
                switch ($state.current) {
                  case 0:
                    $state.returnValue = $this.SomeInt;
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
        }.call(this);
      } else {
        return function (val) {
          var $this = this;
          var $state = {
            current: 0,
            returnValue: null,
          };
          $state.next = function ($callback) {
            try {
              while (true) {
                switch ($state.current) {
                  case 0:
                    $this.SomeInt = val;
                    $state.current = -1;
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
        }.call(this, opt_val);
      }
    };
  });

  $static.AnotherFunction = function (sc) {
    var $state = {
      current: 0,
      returnValue: null,
    };
    var $getValue$1;
    var $setValue$2;
    $state.next = function ($callback) {
      try {
        while (true) {
          switch ($state.current) {
            case 0:
              sc.SomeInt;
              sc.SomeProp().then(function (returnValue) {
                $state.current = 1;
                $getValue$1 = returnValue;
                $state.next($callback);
              }).catch(function (e) {
                $state.error = e;
                $state.current = -1;
                $callback($state);
              });
              return;

            case 1:
              $getValue$1;
              sc.SomeProp(4).then(function (returnValue) {
                $state.current = 2;
                $setValue$2 = returnValue;
                $state.next($callback);
              }).catch(function (e) {
                $state.error = e;
                $state.current = -1;
                $callback($state);
              });
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
