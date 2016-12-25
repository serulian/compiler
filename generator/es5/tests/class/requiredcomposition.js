$module('requiredcomposition', function () {
  var $static = this;
  this.$class('ca41b700', 'First', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function (FirstValue) {
      var instance = new $static();
      instance.FirstValue = FirstValue;
      return instance;
    };
    this.$typesig = function () {
      return {
      };
    };
  });

  this.$class('ea789272', 'Second', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function (SecondValue) {
      var instance = new $static();
      instance.SecondValue = SecondValue;
      return instance;
    };
    this.$typesig = function () {
      return {
      };
    };
  });

  this.$class('4e195caf', 'SomeClass', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function (FirstValue, SecondValue) {
      var instance = new $static();
      instance.First = $g.requiredcomposition.First.new(FirstValue);
      instance.Second = $g.requiredcomposition.Second.new(SecondValue);
      return instance;
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
    this.$typesig = function () {
      return {
      };
    };
  });

  $static.TEST = function () {
    var sc;
    sc = $g.requiredcomposition.SomeClass.new($t.fastbox(42, $g.____testlib.basictypes.Integer), $t.fastbox('hello', $g.____testlib.basictypes.String));
    return $t.fastbox((sc.FirstValue.$wrapped == 42) && $g.____testlib.basictypes.String.$equals(sc.SecondValue, $t.fastbox('hello', $g.____testlib.basictypes.String)).$wrapped, $g.____testlib.basictypes.Boolean);
  };
});
