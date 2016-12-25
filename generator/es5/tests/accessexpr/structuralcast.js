$module('structuralcast', function () {
  var $static = this;
  this.$class('fa55db57', 'BaseClass', true, '', function (T) {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      return instance;
    };
    $instance.Result = $t.property(function () {
      var $this = this;
      return $t.fastbox(true, $g.____testlib.basictypes.Boolean);
    });
    this.$typesig = function () {
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
        "Result|3|43834c3f": true,
      };
      return this.$cachedtypesig = computed;
    };
  });

  this.$class('5a5a0b35', 'SomeClass', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      instance.BaseClass$Integer = $g.structuralcast.BaseClass($g.____testlib.basictypes.Integer).new();
      return instance;
    };
    $instance.Result = $t.property(function () {
      var $this = this;
      return $t.fastbox(false, $g.____testlib.basictypes.Boolean);
    });
    this.$typesig = function () {
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
        "Result|3|43834c3f": true,
      };
      return this.$cachedtypesig = computed;
    };
  });

  $static.DoSomething = function (sc) {
    return;
  };
  $static.TEST = function () {
    var sc;
    sc = $g.structuralcast.SomeClass.new();
    return sc.BaseClass$Integer.Result();
  };
});
