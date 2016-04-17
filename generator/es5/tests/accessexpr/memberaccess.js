$module('memberaccess', function () {
  var $static = this;
  this.$class('SomeClass', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      var init = [];
      init.push($promise.resolve($t.box(2, $g.____testlib.basictypes.Integer)).then(function (result) {
        instance.someInt = result;
      }));
      init.push($promise.resolve($t.box(true, $g.____testlib.basictypes.Boolean)).then(function (result) {
        instance.someBool = result;
      }));
      return $promise.all(init).then(function () {
        return instance;
      });
    };
    $static.Build = function () {
      var $state = $t.sm(function ($continue) {
        while (true) {
          switch ($state.current) {
            case 0:
              $g.memberaccess.SomeClass.new().then(function ($result0) {
                $result = $result0;
                $state.current = 1;
                $continue($state);
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
    $instance.SomeProp = $t.property(function () {
      var $this = this;
      var $state = $t.sm(function ($continue) {
        $state.resolve($this.someInt);
        return;
      });
      return $promise.build($state);
    });
  });

  $static.DoSomething = function (sc, scn) {
    var $state = $t.sm(function ($continue) {
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
              $continue($state);
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
              $continue($state);
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 2:
            $result;
            sc.SomeProp().then(function ($result0) {
              $result = $result0;
              $state.current = 3;
              $continue($state);
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 3:
            $result;
            scn.SomeProp().then(function ($result0) {
              $result = $result0;
              $state.current = 4;
              $continue($state);
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
    var $state = $t.sm(function ($continue) {
      while (true) {
        switch ($state.current) {
          case 0:
            $g.memberaccess.SomeClass.new().then(function ($result0) {
              $result = $result0;
              $state.current = 1;
              $continue($state);
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 1:
            sc = $result;
            $promise.resolve($t.unbox(sc.someBool)).then(function ($result0) {
              $result = $t.box($result0 && $t.unbox(sc.someBool), $g.____testlib.basictypes.Boolean);
              $state.current = 2;
              $continue($state);
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 2:
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
});
