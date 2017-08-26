$module('unwrap', function () {
  var $static = this;
  this.$class('a4aedff5', 'SomeClass', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      return instance;
    };
    $instance.SomeProperty = $t.property(function () {
      var $this = this;
      return $t.fastbox(true, $g.________testlib.basictypes.Boolean);
    });
    this.$typesig = function () {
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
        "SomeProperty|3|f361570c": true,
      };
      return this.$cachedtypesig = computed;
    };
  });

  this.$type('225cfe6c', 'SomeType', false, '', function () {
    var $instance = this.prototype;
    var $static = this;
    this.$box = function ($wrapped) {
      var instance = new this();
      instance[BOXED_DATA_PROPERTY] = $wrapped;
      return instance;
    };
    this.$roottype = function () {
      return $g.unwrap.SomeClass;
    };
    this.$typesig = function () {
      return {
      };
    };
  });

  $static.TEST = function () {
    var st;
    st = $t.fastbox($g.unwrap.SomeClass.new(), $g.unwrap.SomeType);
    return st.$wrapped.SomeProperty();
  };
});
