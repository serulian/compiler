$module('memberaccess', function () {
  var $static = this;
  this.$class('SomeClass', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      var init = [];
      init.push(function () {
        return $promise.wrap(function () {
          $this.someInt = 2;
        });
      }());
      return $promise.all(init).then(function () {
        return instance;
      });
    };
    $static.Build = function () {
      var $returnValue$1;
      var $state = $t.sm(function ($callback) {
        while (true) {
          switch ($state.current) {
            case 0:
              $g.memberaccess.SomeClass.new().then(function (returnValue) {
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
      });
      return $promise.build($state);
    };
    $instance.InstanceFunc = function () {
      return $promise.empty();
    };
    $instance.SomeProp = function () {
      var $this = this;
      var $state = $t.sm(function ($callback) {
        while (true) {
          switch ($state.current) {
            case 0:
              $state.returnValue = $this.someInt;
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

  $static.DoSomething = function (sc, scn) {
    var $returnValue$1;
    var $getValue$2;
    var $getValue$3;
    var $getValue$4;
    var $state = $t.sm(function ($callback) {
      while (true) {
        switch ($state.current) {
          case 0:
            sc.someInt;
            $g.memberaccess.SomeClass.Build;
            $t.dynamicaccess(sc, 'someInt', false, false);
            $g.memberaccess.SomeClass.Build;
            $t.dynamicaccess(scn, 'someInt', false, false);
            $g.maimport.AnotherFunction;
            $g.maimport.AnotherFunction;
            $g.maimport.AnotherFunction;
            sc.InstanceFunc().then(function (returnValue) {
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
            $returnValue$1;
            $t.dynamicaccess(sc, 'InstanceFunc', false, false);
            sc.SomeProp().then(function (returnValue) {
              $state.current = 2;
              $getValue$2 = returnValue;
              $state.next($callback);
            }).catch(function (e) {
              $state.error = e;
              $state.current = -1;
              $callback($state);
            });
            return;

          case 2:
            $getValue$2;
            $t.dynamicaccess(sc, 'SomeProp', true, false)().then(function (returnValue) {
              $state.current = 3;
              $getValue$3 = returnValue;
              $state.next($callback);
            }).catch(function (e) {
              $state.error = e;
              $state.current = -1;
              $callback($state);
            });
            return;

          case 3:
            $getValue$3;
            $t.dynamicaccess(scn, 'SomeProp', true, false)().then(function (returnValue) {
              $state.current = 4;
              $getValue$4 = returnValue;
              $state.next($callback);
            }).catch(function (e) {
              $state.error = e;
              $state.current = -1;
              $callback($state);
            });
            return;

          case 4:
            $getValue$4;
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
