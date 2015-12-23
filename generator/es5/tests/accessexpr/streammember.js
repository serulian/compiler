$module('streammember', function () {
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
  });

  $static.AnotherThing = function (somestream) {
    var $state = {
      current: 0,
      returnValue: null,
    };
    $state.next = function ($callback) {
      try {
        while (true) {
          switch ($state.current) {
            case 0:
              $t.streamaccess(somestream, 'SomeInt');
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
