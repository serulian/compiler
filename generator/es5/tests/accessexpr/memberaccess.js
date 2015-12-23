$module('memberaccess', function () {
  var $static = this;
  this.cls('SomeClass', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function ($callback) {
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
  });

  $static.DoSomething = function (sc, scn) {
    var $state = $t.sm(function ($callback) {
      while (true) {
        switch ($state.current) {
          case 0:
            sc.someInt;
            $g.memberaccess.SomeClass.Build;
            $t.dynamicaccess(sc, 'someInt');
            $g.memberaccess.SomeClass.Build;
            $t.nullaccess(scn, 'someInt');
            $g.maimport.AnotherFunction;
            $g.maimport.AnotherFunction;
            $g.maimport.AnotherFunction;
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
