$module('funcref', function () {
  var $static = this;
  this.$class('SomeClass', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function (value) {
      var instance = new $static();
      var init = [];
      init.push($promise.new(function (resolve) {
        instance.value = value;
        resolve();
      }));
      return $promise.all(init).then(function () {
        return instance;
      });
    };
    $instance.SomeFunction = function () {
      var $this = this;
      var $state = $t.sm(function ($continue) {
        $state.resolve($this.value);
        return;
      });
      return $promise.build($state);
    };
  });

  $static.AnotherFunction = function (toCall) {
    var $state = $t.sm(function ($continue) {
      while (true) {
        switch ($state.current) {
          case 0:
            toCall().then(function ($result0) {
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
  $static.TEST = function () {
    var sc;
    var $state = $t.sm(function ($continue) {
      while (true) {
        switch ($state.current) {
          case 0:
            $g.funcref.SomeClass.new($t.box(true, $g.____testlib.basictypes.Boolean)).then(function ($result0) {
              $temp0 = $result0;
              $result = ($temp0, $temp0);
              $state.current = 1;
              $continue($state);
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 1:
            sc = $result;
            $g.funcref.AnotherFunction($t.dynamicaccess(sc, 'SomeFunction')).then(function ($result0) {
              $result = $result0;
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
