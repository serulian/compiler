$module('requiredcomposition', function () {
  var $static = this;
  this.$class('First', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function (FirstValue) {
      var instance = new $static();
      var init = [];
      init.push($promise.new(function (resolve) {
        instance.FirstValue = FirstValue;
        resolve();
      }));
      return $promise.all(init).then(function () {
        return instance;
      });
    };
  });

  this.$class('Second', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function (SecondValue) {
      var instance = new $static();
      var init = [];
      init.push($promise.new(function (resolve) {
        instance.SecondValue = SecondValue;
        resolve();
      }));
      return $promise.all(init).then(function () {
        return instance;
      });
    };
  });

  this.$class('SomeClass', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function (FirstValue, SecondValue) {
      var instance = new $static();
      var init = [];
      init.push($g.requiredcomposition.First.new(FirstValue).then(function (value) {
        instance.First = value;
      }));
      init.push($g.requiredcomposition.Second.new(SecondValue).then(function (value) {
        instance.Second = value;
      }));
      return $promise.all(init).then(function () {
        return instance;
      });
    };
    Object.defineProperty($instance, 'FirstValue', {
      get: function () {
        return this.First.FirstValue;
      },
      set: function (val) {
        this.First.FirstValue = val;
      },
    });
    Object.defineProperty($instance, 'SecondValue', {
      get: function () {
        return this.Second.SecondValue;
      },
      set: function (val) {
        this.Second.SecondValue = val;
      },
    });
  });

  $static.TEST = function () {
    var sc;
    var $state = $t.sm(function ($callback) {
      while (true) {
        switch ($state.current) {
          case 0:
            $g.requiredcomposition.SomeClass.new($t.nominalwrap(42, $g.____testlib.basictypes.Integer), $t.nominalwrap('hello', $g.____testlib.basictypes.String)).then(function ($result0) {
              $result = $result0;
              $state.current = 1;
              $callback($state);
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 1:
            sc = $result;
            $g.____testlib.basictypes.Integer.$equals(sc.FirstValue, $t.nominalwrap(42, $g.____testlib.basictypes.Integer)).then(function ($result0) {
              return $g.____testlib.basictypes.String.$equals(sc.SecondValue, $t.nominalwrap('hello', $g.____testlib.basictypes.String)).then(function ($result1) {
                $result = $t.nominalwrap($t.nominalunwrap($result0) && $t.nominalunwrap($result1), $g.____testlib.basictypes.Boolean);
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
