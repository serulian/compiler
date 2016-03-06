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
      var $state = $t.sm(function ($callback) {
        while (true) {
          switch ($state.current) {
            case 0:
              $state.resolve($t.nominalwrap(42, $g.____testlib.basictypes.Integer));
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

  this.$interface('Valuable', false, '', function () {
    var $static = this;
  });

  this.$type('Valued', false, '', function () {
    var $instance = this.prototype;
    var $static = this;
    this.new = function ($wrapped) {
      var instance = new this();
      instance.$wrapped = $wrapped;
      return instance;
    };
    this.$box = function (data) {
      var instance = new this();
      instance.$wrapped = data;
      return instance;
    };
    this.$unbox = function (instance) {
      return instance.$wrapped;
    };
    $instance.GetValue = function () {
      var $this = this;
      var $state = $t.sm(function ($callback) {
        while (true) {
          switch ($state.current) {
            case 0:
              $t.nominalunwrap($this).SomeValue().then(function ($result0) {
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
  });

  $static.TEST = function () {
    var sc;
    var v;
    var $state = $t.sm(function ($callback) {
      while (true) {
        switch ($state.current) {
          case 0:
            $g.interface.SomeClass.new().then(function ($result0) {
              $result = $result0;
              $state.current = 1;
              $callback($state);
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 1:
            sc = $result;
            v = sc;
            $t.nominalwrap(v, $g.interface.Valued).GetValue().then(function ($result0) {
              return $g.____testlib.basictypes.Integer.$equals($result0, $t.nominalwrap(42, $g.____testlib.basictypes.Integer)).then(function ($result1) {
                $result = $result1;
                $state.current = 2;
                $callback($state);
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
