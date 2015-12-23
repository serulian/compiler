$module('cast', function () {
  var $static = this;
  this.interface('ISomeInterface', function () {
  });

  this.cls('SomeClass', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function ($callback) {
      var instance = new $static();
      var init = [];
      return $promise.all(init).then(function () {
        return instance;
      });
    };
  });

  $static.DoSomething = function (i) {
    var $state = {
      current: 0,
      returnValue: null,
    };
    $state.next = function ($callback) {
      try {
        while (true) {
          switch ($state.current) {
            case 0:
              $t.cast(i, $g.cast.SomeClass);
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
