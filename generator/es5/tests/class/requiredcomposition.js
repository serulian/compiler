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
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            $g.requiredcomposition.SomeClass.new($t.box(42, $g.____testlib.basictypes.Integer), $t.box('hello', $g.____testlib.basictypes.String)).then(function ($result0) {
              $result = $result0;
              $current = 1;
              $continue($resolve, $reject);
              return;
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 1:
            sc = $result;
            $g.____testlib.basictypes.Integer.$equals(sc.FirstValue, $t.box(42, $g.____testlib.basictypes.Integer)).then(function ($result1) {
              return $promise.resolve($t.unbox($result1)).then(function ($result0) {
                return ($promise.shortcircuit($result0, false) || $g.____testlib.basictypes.String.$equals(sc.SecondValue, $t.box('hello', $g.____testlib.basictypes.String))).then(function ($result2) {
                  $result = $t.box($result0 && $t.unbox($result2), $g.____testlib.basictypes.Boolean);
                  $current = 2;
                  $continue($resolve, $reject);
                  return;
                });
              });
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 2:
            $resolve($result);
            return;

          default:
            $resolve();
            return;
        }
      }
    };
    return $promise.new($continue);
  };
});
