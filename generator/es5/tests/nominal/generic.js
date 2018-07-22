$module('generic', function () {
  var $static = this;
  this.$class('91099e16', 'SomeClass', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      return instance;
    };
    $instance.DoSomething = function () {
      var $this = this;
      return $t.fastbox(true, $g.________testlib.basictypes.Boolean);
    };
    this.$typesig = function () {
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
        "DoSomething|2|cf412abd<aa28dc2d>": true,
      };
      return this.$cachedtypesig = computed;
    };
  });

  this.$type('624a04ee', 'MyType', true, '', function (T) {
    var $instance = this.prototype;
    var $static = this;
    this.$box = function ($wrapped) {
      var instance = new this();
      instance[BOXED_DATA_PROPERTY] = $wrapped;
      return instance;
    };
    this.$roottype = function () {
      return $g.generic.SomeClass;
    };
    $instance.AnotherThing = function () {
      var $this = this;
      return $this.$wrapped.DoSomething();
    };
    $instance.SomeProp = $t.property(function () {
      var $this = this;
      return $t.fastbox(true, $g.________testlib.basictypes.Boolean);
    });
    this.$typesig = function () {
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
        "AnotherThing|2|cf412abd<aa28dc2d>": true,
        "SomeProp|3|aa28dc2d": true,
      };
      return this.$cachedtypesig = computed;
    };
  });

  $static.TEST = function () {
    var m;
    var sc;
    sc = $g.generic.SomeClass.new();
    m = $t.fastbox(sc, $g.generic.MyType($g.________testlib.basictypes.Integer));
    return $t.fastbox(m.SomeProp().$wrapped && m.AnotherThing().$wrapped, $g.________testlib.basictypes.Boolean);
  };
});
