$module('structuralcast', function () {
  var $static = this;
  this.$class('BaseClass', true, '', function (T) {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      var init = [];
      return $promise.all(init).then(function () {
        return instance;
      });
    };
    $instance.Result = $t.property(function () {
      var $this = this;
      var $state = $t.sm(function ($continue) {
        $state.resolve($t.box(true, $g.____testlib.basictypes.Boolean));
        return;
      });
      return $promise.build($state);
    });
  });

  this.$class('SomeClass', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      var init = [];
      init.push($g.structuralcast.BaseClass($g.____testlib.basictypes.Integer).new().then(function (value) {
        instance.BaseClass$Integer = value;
      }));
      return $promise.all(init).then(function () {
        return instance;
      });
    };
    $instance.Result = $t.property(function () {
      var $this = this;
      var $state = $t.sm(function ($continue) {
        $state.resolve($t.box(false, $g.____testlib.basictypes.Boolean));
        return;
      });
      return $promise.build($state);
    });
  });

  $static.DoSomething = function (sc) {
    var $state = $t.sm(function ($continue) {
      sc.BaseClass$Integer;
      $state.resolve();
    });
    return $promise.build($state);
  };
  $static.TEST = function () {
    var sc;
    var $state = $t.sm(function ($continue) {
      while (true) {
        switch ($state.current) {
          case 0:
            $g.structuralcast.SomeClass.new().then(function ($result0) {
              $result = $result0;
              $state.current = 1;
              $continue($state);
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 1:
            sc = $result;
            sc.BaseClass$Integer.Result().then(function ($result0) {
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
