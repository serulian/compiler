$module('basic', function () {
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
    $instance.DoSomething = function () {
      var $this = this;
      var $state = $t.sm(function ($continue) {
        $state.resolve($t.box(true, $g.____testlib.basictypes.Boolean));
        return;
      });
      return $promise.build($state);
    };
  });

  this.$type('MyType', false, '', function () {
    var $instance = this.prototype;
    var $static = this;
    this.$box = function ($wrapped) {
      var instance = new this();
      instance[BOXED_DATA_PROPERTY] = $wrapped;
      return instance;
    };
    $instance.AnotherThing = function () {
      var $this = this;
      var $state = $t.sm(function ($continue) {
        while (true) {
          switch ($state.current) {
            case 0:
              $t.unbox($this).DoSomething().then(function ($result0) {
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
    $instance.SomeProp = $t.property(function () {
      var $this = this;
      var $state = $t.sm(function ($continue) {
        $state.resolve($t.box(true, $g.____testlib.basictypes.Boolean));
        return;
      });
      return $promise.build($state);
    });
  });

  $static.TEST = function () {
    var m;
    var sc;
    var $state = $t.sm(function ($continue) {
      while (true) {
        switch ($state.current) {
          case 0:
            $g.basic.SomeClass.new().then(function ($result0) {
              $result = $result0;
              $state.current = 1;
              $continue($state);
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 1:
            sc = $result;
            m = $t.box(sc, $g.basic.MyType);
            m.SomeProp().then(function ($result1) {
              return $promise.resolve($t.unbox($result1)).then(function ($result0) {
                return ($promise.shortcircuit($result0, false) || m.AnotherThing()).then(function ($result2) {
                  $result = $t.box($result0 && $t.unbox($result2), $g.____testlib.basictypes.Boolean);
                  $state.current = 2;
                  $continue($state);
                });
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
