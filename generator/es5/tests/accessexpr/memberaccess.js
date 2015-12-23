$module('memberaccess', function () {
  var $static = this;
  this.cls('SomeClass', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.$new = function () {
      var instance = new $static();
      function () {
        this.someInt = TODO;
      }.call(instance);
      return instance;
    };
    $static.Build = function () {
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
  $static.DoSomething = function (sc, scn) {
    var $state = {
      current: 0,
      returnValue: null,
    };
    $state.next = function ($callback) {
      try {
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
