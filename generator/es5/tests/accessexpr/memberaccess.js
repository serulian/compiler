$module('memberaccess', function () {
  var $static = this;
  this.$class('SomeClass', false, function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      var init = [];
      init.push($promise.resolve(2).then(function (result) {
        instance.someInt = result;
      }));
      init.push($promise.resolve(true).then(function (result) {
        instance.someBool = result;
      }));
      return $promise.all(init).then(function () {
        return instance;
      });
    };
    $static.Build = function () {
      var $state = $t.sm(function ($callback) {
        while (true) {
          switch ($state.current) {
            case 0:
              $g.memberaccess.SomeClass.new().then(function ($result0) {
                $result = $result0;
                $state.current = 1;
                $callback($state);
              }).catch(function (err) {
                $state.reject(err);
              });
              return;

            case 1:
              $state.resolve($result);
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
      var $this = this;
      return $promise.empty();
    };
    $instance.SomeProp = $t.property(false, function () {
      var $this = this;
      var $state = $t.sm(function ($callback) {
        while (true) {
          switch ($state.current) {
            case 0:
              $state.resolve($this.someInt);
              return;

            default:
              $state.current = -1;
              return;
          }
        }
      });
      return $promise.build($state);
    });
  });

  $static.DoSomething = function (sc, scn) {
    var $state = $t.sm(function ($callback) {
      while (true) {
        switch ($state.current) {
          case 0:
            sc.someInt;
            $t.dynamicaccess($g.memberaccess.SomeClass, 'Build');
            sc.someInt;
            $t.dynamicaccess($g.memberaccess.SomeClass, 'Build');
            scn.someInt;
            $t.dynamicaccess(maimport, 'AnotherFunction');
            $g.maimport.AnotherFunction;
            $g.maimport.AnotherFunction;
            sc.InstanceFunc().then(function ($result0) {
              $result = $result0;
              $state.current = 1;
              $callback($state);
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 1:
            $result;
            $t.dynamicaccess(sc, 'InstanceFunc');
            sc.SomeProp().then(function ($result0) {
              $result = $result0;
              $state.current = 2;
              $callback($state);
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 2:
            $result;
            sc.SomeProp().then(function ($result0) {
              $result = $result0;
              $state.current = 3;
              $callback($state);
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 3:
            $result;
            scn.SomeProp().then(function ($result0) {
              $result = $result0;
              $state.current = 4;
              $callback($state);
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 4:
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
  $static.TEST = function () {
    var sc;
    var $state = $t.sm(function ($callback) {
      while (true) {
        switch ($state.current) {
          case 0:
            $g.memberaccess.SomeClass.new().then(function ($result0) {
              $result = $result0;
              $state.current = 1;
              $callback($state);
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 1:
            sc = $result;
            $state.resolve(sc.someBool && sc.someBool);
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
