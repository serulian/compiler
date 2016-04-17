$module('interface', function () {
  var $static = this;
  this.$class('SomeClass', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      var init = [];
      return $promise.all(init).then(function () {
        return instance;
      });
    };
    $instance.SomeValue = $t.property(function () {
      var $this = this;
      var $state = $t.sm(function ($continue) {
        $state.resolve($t.box(42, $g.____testlib.basictypes.Integer));
        return;
      });
      return $promise.build($state);
    });
  });

  this.$interface('Valuable', false, '', function () {
    var $static = this;
  });

  this.$type('Valued', false, '', function () {
    var $instance = this.prototype;
    var $static = this;
    this.$box = function ($wrapped) {
      var instance = new this();
      instance[BOXED_DATA_PROPERTY] = $wrapped;
      return instance;
    };
    $instance.GetValue = function () {
      var $this = this;
      var $state = $t.sm(function ($continue) {
        while (true) {
          switch ($state.current) {
            case 0:
              $t.unbox($this).SomeValue().then(function ($result0) {
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
  });

  $static.TEST = function () {
    var sc;
    var v;
    var $state = $t.sm(function ($continue) {
      while (true) {
        switch ($state.current) {
          case 0:
            $g.interface.SomeClass.new().then(function ($result0) {
              $result = $result0;
              $state.current = 1;
              $continue($state);
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 1:
            sc = $result;
            v = sc;
            $t.box(v, $g.interface.Valued).GetValue().then(function ($result0) {
              return $g.____testlib.basictypes.Integer.$equals($result0, $t.box(42, $g.____testlib.basictypes.Integer)).then(function ($result1) {
                $result = $result1;
                $state.current = 2;
                $continue($state);
              });
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
