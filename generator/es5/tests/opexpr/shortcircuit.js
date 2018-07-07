$module('shortcircuit', function () {
  var $static = this;
  this.$class('4a28e8dc', 'someError', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      return instance;
    };
    $instance.Message = $t.property(function () {
      var $this = this;
      return $t.fastbox('WHY CALLED? ', $g.________testlib.basictypes.String);
    });
    this.$typesig = function () {
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
        "Message|3|e38ac9b0": true,
      };
      return this.$cachedtypesig = computed;
    };
  });

  $static.neverCalled = function () {
    throw $g.shortcircuit.someError.new();
  };
  $static.anotherNeverCalled = function () {
    throw $g.shortcircuit.someError.new();
  };
  $static.TEST = function () {
    return $t.fastbox(!(false && $g.shortcircuit.neverCalled().$wrapped) || $g.shortcircuit.anotherNeverCalled().$wrapped, $g.________testlib.basictypes.Boolean);
  };
});
