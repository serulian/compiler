$module('property', function () {
  var $static = this;
  this.$class('SomeClass', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      var init = [];
      init.push($promise.resolve($t.box(false, $g.____testlib.basictypes.Boolean)).then(function (result) {
        instance.SomeBool = result;
      }));
      return $promise.all(init).then(function () {
        return instance;
      });
    };
    $instance.SomeProp = $t.property(function () {
      var $this = this;
      var $state = $t.sm(function ($continue) {
        $state.resolve($this.SomeBool);
        return;
      });
      return $promise.build($state);
    }, function (val) {
      var $this = this;
      var $state = $t.sm(function ($continue) {
        $this.SomeBool = val;
        $state.resolve();
      });
      return $promise.build($state);
    });
  });

  $static.AnotherFunction = function (sc) {
    var $state = $t.sm(function ($continue) {
      while (true) {
        switch ($state.current) {
          case 0:
            sc.SomeBool;
            sc.SomeProp().then(function ($result0) {
              $result = $result0;
              $state.current = 1;
              $continue($state);
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 1:
            $result;
            sc.SomeProp($t.box(true, $g.____testlib.basictypes.Boolean)).then(function ($result0) {
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
    var $state = $t.sm(function ($continue) {
      while (true) {
        switch ($state.current) {
          case 0:
            $g.property.SomeClass.new().then(function ($result0) {
              return $g.property.AnotherFunction($result0).then(function ($result1) {
                $result = $result1;
                $state.current = 1;
                $continue($state);
              });
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
});
